import { Component, OnInit } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { GamesService } from "./services/games.service";
import { Game } from "./modules/games";
import { CommonModule } from "@angular/common";
import { MatButtonModule } from '@angular/material/button';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, CommonModule, MatButtonModule],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent implements OnInit{
  title = 'indiestream';
  status: string = '';

  public games: any;

  constructor(private gamesService: GamesService) {  }

  ngOnInit() {
    this.getGames();
  }

  getGames() {
    this.gamesService.getGames().subscribe(
      response => {
        this.games = response;
        //console.log(this.games);
      }
    );
  }

  refreshGame(id: string) {
    this.gamesService.getGame(id).subscribe(
      response => {
        this.games.map((game: Game) => this.games.find((resp: Game) => resp.id === game.id) || game);
      }
    )
  }

  deleteGame(id: string)  {
    this.gamesService.deleteGame(id);
    //TODO refresh?
    //getGames():
  }
}
