import { APP_INITIALIZER, ApplicationConfig, importProvidersFrom } from '@angular/core';
import { provideRouter } from '@angular/router';

import { routes } from './app.routes';
import { provideHttpClient, withInterceptors} from "@angular/common/http";
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { provideAnimations } from "@angular/platform-browser/animations";
import { provideOAuthClient } from "angular-oauth2-oidc";
import { AuthInitializer } from "./services/auth.initializer";
import { AuthInterceptor } from "./services/authInterceptor.service";
import {AuthGuard} from "./guards/auth.guard";

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideHttpClient(withInterceptors([AuthInterceptor])),
    provideAnimationsAsync(),
    provideAnimations(),
    provideOAuthClient(),
    { provide: APP_INITIALIZER,
      useFactory: (authInitializer: AuthInitializer) => () =>
        authInitializer.initializeApp(),
      deps: [AuthInitializer],
      multi: true, },
    AuthGuard,
  ]
};
