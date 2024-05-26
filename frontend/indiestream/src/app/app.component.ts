import { Component, OnInit } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { GamesService } from "./services/games.service";
import { Game } from "./modules/games";
import { CommonModule } from "@angular/common";
import { MatButtonModule } from '@angular/material/button';
import { HttpClientModule } from "@angular/common/http";
import { GameUploadComponent} from "./components/game-upload/game-upload.component";
import {GamesOverviewComponent} from "./components/games-overview/games-overview.component";

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, CommonModule, MatButtonModule, HttpClientModule, GameUploadComponent, GamesOverviewComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent {
  title = 'indiestream';
  status: string = '';

  constructor() {  }
}
