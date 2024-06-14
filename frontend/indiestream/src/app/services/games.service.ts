import { Injectable } from "@angular/core";
import { HttpClient, HttpEvent, HttpHeaders} from '@angular/common/http';
import { Games, Game } from '../modules/games';
import { Observable } from "rxjs";
import { AppConfigService } from "./app-config.service";
import { AuthService} from "./auth.service";
import {FormGroup} from "@angular/forms";

@Injectable({
  providedIn: 'root'
})
export class GamesService {
  private apiUrl = this.configService.getConfig().apiUrl;
  constructor(private http: HttpClient, private configService: AppConfigService) {
    console.log(this.configService.getConfig())
  }

  getGames(): Observable<Games> {
    const httpOptions = {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    };
    return this.http.get<Games>(this.apiUrl + "/games");
  }

  getGame(id: string): Observable<Game> {
    return this.http.get<Game>(this.apiUrl + "/games/" + id)
  }

  deleteGame(id: string): void {
    this.http.delete(this.apiUrl + "/games/" + id)
  }

  uploadGame(gameForm: FormGroup): Observable<HttpEvent<Object>> {
    const formData = new FormData();
    formData.append('title', gameForm.get('title')?.value);
    const files: FileList = gameForm.get('file')?.value;
    if (files && files.length > 0) {
      formData.append('file', files[0]);
    }
    return this.http.post(this.apiUrl + "/games", formData, {
      reportProgress: true,
      observe: 'events'
    });
  }
}

