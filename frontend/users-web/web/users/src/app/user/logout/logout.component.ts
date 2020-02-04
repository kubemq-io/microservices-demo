import { Component, OnInit } from '@angular/core';
import {UsersApiService} from "../../shared/users-api.service";
import {ToastrService} from "ngx-toastr";
import {NgForm} from "@angular/forms";
import {ApiResponse} from "../../shared/user.model";

@Component({
  selector: 'app-logout',
  templateUrl: './logout.component.html',
  styleUrls: ['./logout.component.css']
})
export class LogoutComponent implements OnInit {

  constructor(public apiService: UsersApiService, private toastr: ToastrService) {}

  ngOnInit() {
    this.resetForm()
  }
  onSubmit(form: NgForm) {

    this.apiService.postLogout().subscribe((response:ApiResponse)=> {
      if (response.is_error) {
        this.toastr.error(response.message,'Logout');
        return
      } else {
        this.toastr.success("user logged out successfully",'Logout');
      }
    },
      error =>   this.toastr.error('Connection Error','Logout'))
  }

  resetForm(form?: NgForm) {
    if (form != null)
      form.resetForm();

    this.apiService.logoutFormData = {
      user_id:'',
    }
  }
}
