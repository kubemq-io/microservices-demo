import { Component, OnInit } from '@angular/core';
import { UsersApiService } from '../../shared/users-api.service';
import { NgForm } from '@angular/forms';
import {ApiResponse} from '../../shared/user.model';
import { ToastrService } from 'ngx-toastr';

@Component({
  selector: 'app-verify',
  templateUrl: './verify.component.html',
  styleUrls: ['./verify.component.css']
})
export class VerifyComponent implements OnInit {
  constructor(public apiService: UsersApiService, private toastr: ToastrService) { }
  ngOnInit() {
    this.resetForm()
  }

  onSubmit(form: NgForm) {

    this.apiService.postVerify().subscribe((response: ApiResponse)=> {
      if (response.is_error) {
        this.toastr.error(response.message, 'Verify');
        return
      } else {
        this.toastr.success("user verified successfully", 'Verify');
      }
    },
      error =>   this.toastr.error(`Connection Error`, 'Verify'))
  }

  resetForm(form?: NgForm) {
    if (form != null)
      form.resetForm();

    this.apiService.verifyFormData = {
      name: '',
      token:'',
    }
  }

}
