import {Component, OnInit} from '@angular/core';
import {MatButton, MatFabButton} from "@angular/material/button";
import {NgForOf, NgIf} from "@angular/common";
import {GamesService} from "../../services/games.service";
import {Game} from "../../modules/games";
import {
  MatCell, MatCellDef,
  MatColumnDef,
  MatHeaderCell, MatHeaderCellDef,
  MatHeaderRow,
  MatHeaderRowDef,
  MatRow,
  MatRowDef,
  MatTable
} from "@angular/material/table";
import {HttpClientModule} from "@angular/common/http";
import {AuthService} from "../../services/auth.service";
import {MatIcon} from "@angular/material/icon";
import {RouterLink} from "@angular/router";
import {
  MatCard,
  MatCardActions,
  MatCardContent,
  MatCardHeader,
  MatCardSubtitle,
  MatCardTitle
} from "@angular/material/card";
import {OAuthService} from "angular-oauth2-oidc";
import {ClipboardModule} from "@angular/cdk/clipboard";

@Component({
  selector: 'app-games-overview',
  standalone: true,
  imports: [
    ClipboardModule,
    MatButton,
    NgForOf,
    MatTable,
    MatColumnDef,
    MatHeaderCell,
    MatCell,
    MatHeaderRow,
    MatRow,
    MatRowDef,
    MatHeaderRowDef,
    MatCellDef,
    MatHeaderCellDef,
    HttpClientModule,
    NgIf,
    MatIcon,
    RouterLink,
    MatCard,
    MatCardTitle,
    MatCardSubtitle,
    MatCardHeader,
    MatCardContent,
    MatCardActions,
    MatFabButton,
  ],
  templateUrl: './games-overview.component.html',
  styleUrl: './games-overview.component.scss'
})
export class GamesOverviewComponent implements OnInit {
  public games: Game[] = [];

  constructor(private gamesService: GamesService, public authService: AuthService, private oAuthService: OAuthService) {

  }

  ngOnInit() {
    this.getGames();
  }

  getGames() {
    this.gamesService.getGames().subscribe(
      response => {
        this.games = response;
        this.games.filter((game: Game) => !(game.url)).forEach((game: Game) => {
          this.refreshGame(game.id);
        });
      }
    );
  }

  refreshGame(id: string) {
    this.gamesService.getGame(id).subscribe(
      (response: Game) => {
        this.games = this.games.map((game: Game) => game.id === id ? response : game);
      }
    );
  }

  deleteGame(id: string) {
    this.gamesService.deleteGame(id).subscribe(
      response => {
        if (response.status === 204) {
          this.games = this.games.filter(game => game.id !== id);
        }
      }
    );
  }
}
