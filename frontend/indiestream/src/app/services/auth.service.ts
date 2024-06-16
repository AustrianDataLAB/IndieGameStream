import { Injectable } from '@angular/core';
import { OAuthService } from 'angular-oauth2-oidc';
import {Router} from "@angular/router";

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  constructor( private oAuthService: OAuthService, private router: Router ) {}

  getName(): string {
    const claims = this.oAuthService.getIdentityClaims();
    return claims['name'];
  }

  getEmail(): string {
    const claims = this.oAuthService.getIdentityClaims();
    return claims['email'];
  }

  getPictureUrl(): string {
    const claims = this.oAuthService.getIdentityClaims();
    if (claims === null)  {
      return '';
    }
    return claims['picture'];
  }

  getIdToken(): string {
    return this.oAuthService.getIdToken();
  }

  isAuthenticated(): boolean {
    return (
      this.oAuthService.hasValidIdToken() &&
      this.oAuthService.hasValidAccessToken()
    );
  }

  login() {
    this.oAuthService.initCodeFlow(this.router.url);
  }

  logout() {
    this.oAuthService.revokeTokenAndLogout();
  }
}
