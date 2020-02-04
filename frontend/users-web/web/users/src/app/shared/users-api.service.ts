import { Injectable } from '@angular/core';
import {LoginRequest, LogoutRequest, RegisterRequest, VerifyRequest} from './user.model';
import {HttpClient} from '@angular/common/http';
import { EnvService } from '../env.service';

@Injectable({
  providedIn: 'root'
})
export class UsersApiService {
  registerFormData: RegisterRequest;
  verifyFormData: VerifyRequest;
  loginFormData: LoginRequest;
  logoutFormData: LogoutRequest;

  constructor(private http: HttpClient, private env: EnvService) { }


  postRegister() {
    return this.http.post(this.env.apiUrl + '/register', this.registerFormData);
  }

  postVerify() {
    return this.http.post(this.env.apiUrl + '/verify', this.verifyFormData);
  }

  postLogin() {
    return this.http.post(this.env.apiUrl + '/login', this.loginFormData);
  }

  postLogout() {
    return this.http.post(this.env.apiUrl + '/logout', this.logoutFormData);
  }
}
