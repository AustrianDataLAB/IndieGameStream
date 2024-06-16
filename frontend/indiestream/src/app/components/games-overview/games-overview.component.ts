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

@Component({
    selector: 'app-games-overview',
    standalone: true,
  imports: [
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
export class GamesOverviewComponent implements OnInit
{
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

  refreshAllGames() {
      for (let game of this.games) {
        this.refreshGame(game.id);
      }
  }

    deleteGame(id: string) {
        this.gamesService.deleteGame(id);
        //TODO refresh?
        //getGames():
    }
}
