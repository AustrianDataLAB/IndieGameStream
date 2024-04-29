import { Injectable } from "@angular/core";
import { HttpClient } from '@angular/common/http';
import { Games } from './games';
import { Observable } from "rxjs";
import { environment } from "./environment";

@Injectable({
  providedIn: 'root'
})
export class GamesService {
  private apiUrl = environment.apiUrl
  constructor(private http: HttpClient) {  }

  getGames(): Observable<Games> {
    return this.http.get<Games>(this.apiUrl + "/games/");
  }
}

