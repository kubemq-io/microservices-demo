import { Component, OnInit } from '@angular/core';
import { UsersApiService } from '../shared/users-api.service';

@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.css']
})
export class UserComponent implements OnInit {

  constructor(private apiService: UsersApiService) { }

  ngOnInit() {
  }

}
