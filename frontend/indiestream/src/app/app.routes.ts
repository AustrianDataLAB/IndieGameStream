import { Route } from '@angular/router';
import {GameUploadComponent} from "./components/game-upload/game-upload.component";
import {GamesOverviewComponent} from "./components/games-overview/games-overview.component";
import {AccountComponent} from "./components/account/account.component";
import {LandingPageComponent} from "./components/landing-page/landing-page.component";
import {AuthGuard} from "./guards/auth.guard";
import {AppComponent} from "./app.component";
import {LayoutComponent} from "./components/layout/layout.component";

export const routes: Route[] = [
  {
    path: '',
    component: LandingPageComponent,
  },
  {
    path: '',
    component: LayoutComponent,
    canActivate: [AuthGuard],
    children: [
      { path: 'dashboard', component: GamesOverviewComponent },
      { path: 'upload', component: GameUploadComponent },
      { path: 'account', component: AccountComponent },
      { path: '**', redirectTo: 'dashboard', pathMatch: 'full' }
    ]
  },
];
