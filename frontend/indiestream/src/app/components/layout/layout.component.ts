import { Component, ViewChild } from '@angular/core';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { CommonModule, NgOptimizedImage } from "@angular/common";
import { MatButtonModule } from '@angular/material/button';
import { HttpClientModule } from "@angular/common/http";
import { MatToolbar } from "@angular/material/toolbar";
import { MatIcon } from "@angular/material/icon";
import { MatSidenav, MatSidenavContainer, MatSidenavContent } from '@angular/material/sidenav';
import { MatListItem, MatNavList } from "@angular/material/list";
import { OAuthService } from "angular-oauth2-oidc";
import { GameUploadComponent } from "../game-upload/game-upload.component";
import { GamesOverviewComponent } from "../games-overview/games-overview.component";
import { AuthService } from "../../services/auth.service";

@Component({
  selector: 'app-layout',
  standalone: true,
  imports: [
    RouterOutlet,
    CommonModule,
    MatButtonModule,
    HttpClientModule,
    GameUploadComponent,
    GamesOverviewComponent,
    MatToolbar,
    MatIcon,
    MatSidenav,
    MatSidenavContainer,
    MatNavList,
    MatListItem,
    MatSidenavContent,
    MatSidenav,
    MatSidenavContainer,
    RouterLink,
    RouterLinkActive,
    NgOptimizedImage
  ],
  templateUrl: './layout.component.html',
  styleUrl: './layout.component.scss'
})
export class LayoutComponent {
  title = 'IndieGameStream';
  status: string = '';
  @ViewChild(MatSidenav)
  sidenav!: MatSidenav;
  uncollapsed = true;

  constructor(public authService: AuthService, private oAuthService: OAuthService) { }

  toggleMenu() {
    this.uncollapsed = !this.uncollapsed;
  }

  logout() {
    this.authService.logout();
  }
}
