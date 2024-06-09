import {inject} from '@angular/core';
import {
  HttpInterceptorFn,
} from '@angular/common/http';
import { AuthService } from "./auth.service";

export const AuthInterceptor: HttpInterceptorFn = (req, next) => {
  const authService = inject(AuthService);
  const idToken = authService.getIdToken();

  // bearer token is not allowed to be sent to this url
  if (idToken && req.url !== 'https://accounts.google.com/.well-known/openid-configuration') {
    const cloned = req.clone({
      headers: req.headers.set('Authorization', `Bearer ${idToken}`)
    });
    return next(cloned);
  } else {
    return next(req);
  }
};
