import { Component, OnInit } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { GamesService } from "./games.service";
import {Games} from "./games";

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent implements OnInit{
  title = 'indiestream';
  status: string = '';

  public games: any;
  public pong: any;

  constructor(private gamesService: GamesService) {  }

  ngOnInit() {
    this.getGames();
  }

  getGames() {
    this.gamesService.getGames().subscribe(
      response => {
        this.games = response;
        console.log(this.games);
      }
    );
  }
}
