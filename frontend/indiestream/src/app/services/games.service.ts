import { Injectable } from "@angular/core";
import {HttpClient, HttpEvent} from '@angular/common/http';
import { Games, Game } from '../modules/games';
import { Observable } from "rxjs";
import { AppConfigService } from "./app-config.service";

@Injectable({
  providedIn: 'root'
})
export class GamesService {
  private apiUrl = this.configService.getConfig().apiUrl;
  constructor(private http: HttpClient, private configService: AppConfigService) {
    console.log(this.configService.getConfig())
  }

  getGames(): Observable<Games> {
    return this.http.get<Games>(this.apiUrl + "/games/");
  }

  getGame(id: string): Observable<Game> {
    return this.http.get<Game>(this.apiUrl + "/games/" + id + "/")
  }

  deleteGame(id: string): void{
    this.http.delete(this.apiUrl + "/games/" + id + "/")
  }

  uploadGame(file: File): Observable<HttpEvent<Object>>{
    return this.http.post(this.apiUrl + "/games/", file, {
      reportProgress: true,
      observe: 'events'
    });
  }
}

