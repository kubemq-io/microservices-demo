export interface ApiResponse {
  is_error: boolean;
  message: string;
  data: object;
}


export interface User {
    id: string;
    create_at: string;
    name: string;
    password: string;
    email: string;
    state: number;
    token: string;
}

export class RegisterRequest {
    name: string;
    password: string;
    email: string;
}


export class VerifyRequest {
  name: string;
  token: string;

}


export class LoginRequest {
  name: string;
  password: string;

}

export class LogoutRequest {
  user_id:string

}
