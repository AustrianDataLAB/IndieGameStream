import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { firstValueFrom, lastValueFrom } from 'rxjs';

export interface AppConfig {
  production: boolean;
  apiUrl: string;
}

@Injectable({
  providedIn: 'root'
})
export class AppConfigService {
  private config: AppConfig = {
    production: false,
    apiUrl: "http://localhost:8080"
  };

  constructor(private http: HttpClient) {}

  loadConfig(): Promise<void> {
    return firstValueFrom(
      this.http.get<AppConfig>('/assets/app.config.json')
    ).then( data => {
      console.log("Config file loaded successfully", data)
      this.config = data;
    }).catch( err => {
      console.error(err);
      return Promise.reject(err)
    });
  }

  getConfig(): AppConfig {
    return this.config;
  }
}
