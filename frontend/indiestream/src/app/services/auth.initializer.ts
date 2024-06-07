import { AuthConfig, OAuthService } from "angular-oauth2-oidc";
import { Injectable, isDevMode } from '@angular/core';
import { Router } from "@angular/router";
import { AuthService } from "./auth.service";

@Injectable({
  providedIn: 'root',
})
export class AuthInitializer {

  constructor( private oauthService: OAuthService, private router: Router, private authService: AuthService ) {
  }

  public initializeApp(): Promise<void> {
    return new Promise((resolve, reject) => {
      const authConfig: AuthConfig = {
        issuer: 'https://accounts.google.com',
        strictDiscoveryDocumentValidation: false,
        clientId: '516825360638-s63uq4ecthilcghh1r1ojbhuqlo3s2ef.apps.googleusercontent.com',
        dummyClientSecret: 'redacted',
        redirectUri: window.location.origin,
        logoutUrl: window.location.origin,
        oidc: true,
        scope: 'openid profile email',
        responseType: 'code',
        showDebugInformation: isDevMode(),
      };
      this.oauthService.configure(authConfig);
      this.oauthService.setupAutomaticSilentRefresh();
      this.oauthService.loadDiscoveryDocumentAndLogin().then(() => {
        /*
        if (this.oauthService.hasValidIdToken() && this.oauthService.hasValidAccessToken()) {
          const url = decodeURIComponent(<string>this.oauthService.state);
          this.router.navigateByUrl(url);
        }*/
        resolve();
      }).catch(err => {
        console.error('OAuth initialization error', err);
        reject(err);
      });
    });
  }
}
