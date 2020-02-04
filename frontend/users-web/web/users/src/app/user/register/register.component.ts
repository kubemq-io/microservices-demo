import { Component, OnInit } from '@angular/core';
import { UsersApiService } from '../../shared/users-api.service';
import { NgForm } from '@angular/forms';
import {RegisterRequest, User} from "../../shared/user.model";
import {ApiResponse} from "../../shared/user.model";
import { ToastrService } from 'ngx-toastr';
@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']
})
export class RegisterComponent implements OnInit {
  constructor(public apiService: UsersApiService, private toastr: ToastrService) { }
  user: User;
  ngOnInit() {
    this.resetForm()
  }

  onSubmit(form: NgForm) {

    this.apiService.postRegister().subscribe((response:ApiResponse)=> {
      if (response.is_error) {
        this.toastr.error(response.message,'Register');
        return
      } else {
         this.user = response.data as User;

        this.apiService.verifyFormData.name=this.user.name;
        this.apiService.verifyFormData.token=this.user.token;

        this.apiService.loginFormData.name=this.apiService.registerFormData.name;
        this.apiService.loginFormData.password=this.apiService.registerFormData.password;

        this.apiService.logoutFormData.user_id=this.user.id;
        this.toastr.success(this.user.name + ' registered successfully','Register');
      }
    },
      error =>   this.toastr.error('Connection Error','Register'))
  }

  resetForm(form?: NgForm) {
    if (form != null)
      form.resetForm();

    this.apiService.registerFormData = {
      name: '',
      password:'',
      email: '',
    }
  }

}
