import { AuthConfig, OAuthService } from "angular-oauth2-oidc";
import { Injectable, isDevMode } from '@angular/core';
import { Router } from "@angular/router";
import { AuthService } from "./auth.service";
import {AppConfigService} from "./app-config.service";

@Injectable({
  providedIn: 'root',
})
export class AuthInitializer {

  constructor( private oauthService: OAuthService, private appConfigService: AppConfigService,
               private router: Router ) {
  }

  public initializeApp(): Promise<void> {
    this.appConfigService.loadConfig();
    return new Promise((resolve, reject) => {
      const authConfig: AuthConfig = {
        issuer: 'https://accounts.google.com',
        strictDiscoveryDocumentValidation: false,
        clientId: '516825360638-ai7mibm97c1i5o66l18iqlfuqffl1dba.apps.googleusercontent.com',
        redirectUri: window.location.origin,
        logoutUrl: window.location.origin,
        oidc: true,
        scope: 'openid profile email',
        showDebugInformation: isDevMode(),
      };
      this.oauthService.configure(authConfig);
      this.oauthService.setupAutomaticSilentRefresh();
      this.oauthService.loadDiscoveryDocumentAndTryLogin().then(() => {
        if (this.oauthService.hasValidIdToken() && this.oauthService.hasValidAccessToken()) {
          const url = decodeURIComponent(<string>this.oauthService.state);
          this.router.navigateByUrl(url);
        }
        resolve();
      }).catch(err => {
        console.error('OAuth initialization error', err);
        reject(err);
      });
    });
  }
}
