import { APP_INITIALIZER, ApplicationConfig, importProvidersFrom } from '@angular/core';
import { provideRouter, Router } from '@angular/router';

import { routes } from './app.routes';
import { HTTP_INTERCEPTORS, provideHttpClient } from "@angular/common/http";
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { provideAnimations } from "@angular/platform-browser/animations";
import { provideOAuthClient } from "angular-oauth2-oidc";
import { AuthInitializer } from "./services/auth.initializer";
import { AuthInterceptor} from "./services/authInterceptor.service";
import { AppConfigService } from './services/app-config.service';

export function initConfig(appConfig: AppConfigService) {
  return () => appConfig.loadConfig();
}

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideHttpClient(),
    provideAnimationsAsync(),
    provideAnimations(),
    provideOAuthClient(),
    { provide: APP_INITIALIZER,
      useFactory: (authInitializer: AuthInitializer) => () =>
        authInitializer.initializeApp(),
      deps: [AuthInitializer],
      multi: true, },
    {
      provide: HTTP_INTERCEPTORS,
      useClass: AuthInterceptor,
      multi: true,
    },]
};
