import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { map } from 'rxjs/operators';


@Injectable({
  providedIn: 'root'
})
export class ApiService {

  signinUrl = 'api/user/signin'
  // signoutUrl = 'api/user/signout'

  constructor(private http:HttpClient, private router: Router) { }

  signin(loginData: {email: string; password: string;}) {
    // var res : SellerLoginResponse;
    return this.http.post<any>(this.signinUrl, loginData).pipe(map(
      data => {
        return data
        // return res
      }
    ))
  }

  // signout() {
  //   return this.http.get<any>(this.signoutUrl)
  // }
}
