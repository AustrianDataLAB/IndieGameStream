import {Component, OnInit} from '@angular/core';
import {MatButton} from "@angular/material/button";
import {NgForOf} from "@angular/common";
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
    ],
    templateUrl: './games-overview.component.html',
    styleUrl: './games-overview.component.scss'
})
export class GamesOverviewComponent implements OnInit
{
    columnsToDisplay = ['ID', 'title', 'status', 'url', 'refresh', 'delete'];
    public games: any;

    constructor(private gamesService: GamesService) {

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

    deleteGame(id: string) {
        this.gamesService.deleteGame(id);
        //TODO refresh?
        //getGames():
    }

}
