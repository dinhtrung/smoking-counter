import { Component, OnInit, OnDestroy } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

import SharedModule from 'app/shared/shared.module';
import { AccountService } from 'app/core/auth/account.service';
import { Account } from 'app/core/auth/account.model';
import {HttpClient} from "@angular/common/http";
import {ApplicationConfigService} from "../core/config/application-config.service";
import {Smoke} from "./smoke.model";
import dayjs from "dayjs/esm";
import {success} from "concurrently/dist/src/defaults";

@Component({
  standalone: true,
  selector: 'jhi-home',
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss',
  imports: [SharedModule, RouterModule],
})
export default class HomeComponent implements OnInit, OnDestroy {
  account: Account | null = null;
  smokes: Smoke[] = [];

  private readonly destroy$ = new Subject<void>();

  constructor(
    private accountService: AccountService,
    private httpClient: HttpClient,
    private applicationConfigService: ApplicationConfigService,
    private router: Router,
  ) {}

  ngOnInit(): void {
    this.accountService
      .getAuthenticationState()
      .pipe(takeUntil(this.destroy$))
      .subscribe(account => {
        this.account = account;
        if (account) {
          this.getAllEvents(account)
        }
      });
  }

  login(): void {
    this.router.navigate(['/login']);
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  getAllEvents(account: Account): void {
    this.httpClient.get<Smoke[]>(this.applicationConfigService.getEndpointFor('api/smokes'))
   .subscribe(smokes => this.smokes = smokes);
  }

  addSmoke(): void {
    const hour = dayjs().format('HH:mm');
    this.httpClient.post(this.applicationConfigService.getEndpointFor('api/smokes'), hour)

      .subscribe();
  }
}
