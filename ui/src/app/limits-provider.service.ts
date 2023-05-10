import {Injectable} from '@angular/core';
import {Limits } from "./limits";
import {HttpClient} from "@angular/common/http";

@Injectable({
  providedIn: 'root'
})
export class LimitsProviderService {

  limits: Limits = {max_batch_size: 50, max_nodes: 100}
  private interval: any = null

  private limitsEndpoint = "/api/v1/limits"

  constructor(private httpClient: HttpClient) {
    this.updateLimits()
    this.interval = setInterval(() => {
      this.updateLimits()
    }, 30 * 60 * 1000)
  }

  get maxNodes() {
    return this.limits.max_nodes
  }

  get maxBatchSize() {
    return this.limits.max_batch_size
  }

  updateLimits(): void {
    this.httpClient.get<Limits>(this.limitsEndpoint).subscribe(value => {
      this.limits = value
    })
  }
}
