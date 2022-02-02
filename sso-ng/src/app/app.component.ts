import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ApiService } from './api.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit{
  title = 'sso';

  loginForm = this.fb.group({
    email: ['', Validators.email],
    password: ['', Validators.required]
  });
  
  constructor(private api: ApiService, private fb: FormBuilder) {}

  onSubmit(): void{
    if (this.loginForm.valid){
      console.log(this.loginForm.value)
      this.api.signin(this.loginForm.value).subscribe((data: any) => {
        console.log(data)
        if (data.seller != null) {
          localStorage.setItem('user', JSON.stringify(data.seller))
          // this.router.navigate(['seller/dashboard'])
        }
        // if (data.access_token && data.refresh_token != null) {
        //   localStorage.setItem('seller_access_token', data.access_token)
        //   localStorage.setItem('seller_refresh_token', data.refresh_token)
        //   this.router.navigate(['seller/dashboard'])
        // }
      })
    }
  }

  ngOnInit(): void {
      // this.api.signin({email: "mail@amit.kr", password: "Amit123"}).subscribe()
      // this.api.signin1({email: "mail@amit.kr", password: "Amit123"}).subscribe()
      // this.api.signin2({email: "mail@amit.kr", password: "Amit123"}).subscribe()
  }
}


