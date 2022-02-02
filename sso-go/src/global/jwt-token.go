package global

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"time"

	"github.com/amitdotkr/sso/sso-go/src/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ExtractToken(c *fiber.Ctx) string {
	bearToken := c.Cookies("access_token")
	return bearToken
}

func TokenClaims(token string) (jwt.MapClaims, error) {
	parsedToken, err := ParsingToken(token)
	if err != nil {
		return nil, err
	}
	claims := parsedToken.Claims.(jwt.MapClaims)
	return claims, nil
}

func ParsingToken(token string) (*jwt.Token, error) {
	pubKey, err := ioutil.ReadFile(PUBKEY_LOC)
	if err != nil {
		log.Fatalln(err)
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(pubKey)
	if err != nil {
		log.Printf("error: %v", err.Error())
	}

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			alg := token.Header["alg"]
			log.Printf("unexpected signing method: %v", alg)
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	return jwtToken, nil
}

// Regenerate Refresh Token Using Refresh Token
// Required
func RegenerateRefreshToken(c *fiber.Ctx, userid, role, keyId string) error {

	kid, err := primitive.ObjectIDFromHex(keyId)
	if err != nil {
		return err
	}

	prvKey, err := ioutil.ReadFile(PRVKEY_LOC)
	if err != nil {
		return err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	if err != nil {
		return err
	}

	exp := time.Now().Add(time.Hour * time.Duration(RefreshTokenTime))
	tokenExpire := exp.Unix()

	rtClaims := jwt.MapClaims{}
	rtClaims["kid"] = keyId
	rtClaims["uid"] = userid
	rtClaims["role"] = role
	rtClaims["type"] = "refresh"
	rtClaims["exp"] = tokenExpire

	rt, err := jwt.NewWithClaims(jwt.SigningMethodRS256, rtClaims).SignedString(key)
	if err != nil {
		return err
	}

	if err := RefreshTokenInDatabase(userid, tokenExpire, kid, rt); err != nil {
		return err
	}

	rc := new(fiber.Cookie)
	rc.Name = "refresh_token"
	rc.Value = rt
	rc.Expires = exp
	rc.HTTPOnly = true
	rc.Secure = CookieSecure
	rc.SameSite = "Strict"

	c.Cookie(rc)

	return nil
}

// Saving Refresh Token in Database
// Required
func RefreshTokenInDatabase(userid string, tokenExpire int64, kid primitive.ObjectID, rt string) error {
	uid, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return err
	}

	t := time.Unix(tokenExpire, 0)

	filter := bson.M{"_id": kid}

	update := bson.M{
		"$set": bson.M{
			"token":    rt,
			"uid":      uid,
			"expireAt": t,
		},
	}
	opts := options.Update().SetUpsert(true)
	RefreshCollection.UpdateOne(Ctx, filter, update, opts)

	return nil
}

// Required
func CreateAccessTokenGo(c *fiber.Ctx, userid string, role string) error {

	prvKey, err := ioutil.ReadFile(PRVKEY_LOC)
	if err != nil {
		log.Printf("err: %v", err.Error())
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	if err != nil {
		return err
	}

	exp := time.Now().Add(time.Minute * time.Duration(AccessTokenTime))

	atClaims := jwt.MapClaims{}
	atClaims["uid"] = userid
	atClaims["role"] = role
	atClaims["type"] = "access"
	atClaims["exp"] = exp.Unix()

	at, err := jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims).SignedString(key)
	if err != nil {
		return err
	}

	ac := new(fiber.Cookie)
	ac.Name = "access_token"
	ac.Value = at
	ac.Expires = exp
	ac.HTTPOnly = true
	ac.Secure = CookieSecure
	ac.SameSite = "Strict"

	// log.Printf("Secure Val: %v", ac.Secure)

	c.Cookie(ac)

	return nil
}

// Required
func CreateRefreshTokenGo(c *fiber.Ctx, userid string, role string) error {

	kid := primitive.NewObjectID()
	keyId := kid.Hex()

	prvKey, err := ioutil.ReadFile(PRVKEY_LOC)
	if err != nil {
		return err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	if err != nil {
		return err
	}

	exp := time.Now().Add(time.Hour * time.Duration(RefreshTokenTime))
	tokenExpire := exp.Unix()

	rtClaims := jwt.MapClaims{}
	rtClaims["kid"] = keyId
	rtClaims["uid"] = userid
	rtClaims["role"] = role
	rtClaims["type"] = "refresh"
	rtClaims["exp"] = tokenExpire

	rt, err := jwt.NewWithClaims(jwt.SigningMethodRS256, rtClaims).SignedString(key)
	if err != nil {
		return err
	}

	rc := new(fiber.Cookie)
	rc.Name = "refresh_token"
	rc.Value = rt
	rc.Expires = exp
	rc.HTTPOnly = true
	rc.Secure = CookieSecure
	rc.SameSite = "Strict"

	c.Cookie(rc)

	if err := RefreshTokenInDatabase(userid, tokenExpire, kid, rt); err != nil {
		return err
	}

	return nil
}

func RegenerateTokenUsingRefreshToken(c *fiber.Ctx) error {
	token := c.Cookies("refresh_token")

	claims, err := TokenClaims(token)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": entities.Error{Type: "Authentication Error", Detail: err.Error()},
		})
	}

	claimType := claims["type"].(string)
	if claimType != "refresh" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": entities.Error{Type: "Authentication Error:", Detail: "Token is not access type"},
		})
	}

	userId := claims["uid"].(string)
	keyId := claims["kid"].(string)

	kid, err := primitive.ObjectIDFromHex(keyId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": entities.Error{Type: "Object Id Error", Detail: err.Error()},
		})
	}

	data := &entities.Refreshtoken{}
	filter := bson.M{"_id": kid}
	if err := RefreshCollection.FindOne(Ctx, filter).Decode(data); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": entities.Error{Type: "Not Found", Detail: err.Error()},
		})
	}

	if token != data.Token {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": entities.Error{Type: "Token Not Found", Detail: "Refresh Token is not available in DB"},
		})
	}

	if err := CreateTokenPairUsingRefreshToken(c, userId, "seller", keyId); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error: ": entities.Error{
				Type:   "Token Parsing Error",
				Detail: err.Error()},
		})
	}
	return nil
}

func CreateTokenPairUsingRefreshToken(c *fiber.Ctx, userid, role, keyId string) error {
	if err := CreateAccessTokenGo(c, userid, role); err != nil {
		return err
	}
	if err := RegenerateRefreshToken(c, userid, role, keyId); err != nil {
		return err
	}

	return nil
}

// Required
func ValidatingUser(c *fiber.Ctx) (string, error) {
	var userId string

	at := c.Cookies("access_token", "")

	if at != "" {
		claims, err := TokenClaims(at)
		if err != nil {
			return "", err
		}
		userId = claims["uid"].(string)
	}

	if at == "" {
		rt := c.Cookies("refresh_token", "")
		if rt == "" {
			return "", errors.New("Refresh Token is Unavailable/Expired")
		}
		RegenerateTokenUsingRefreshToken(c)

		claims, err := TokenClaims(rt)
		if err != nil {
			return "", err
		}
		userId = claims["uid"].(string)
	}
	return userId, nil
}
