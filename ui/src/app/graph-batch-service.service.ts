import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Observable} from "rxjs";
import {GraphBatchList} from "./graph-batch-list";
import {GraphBatch} from "./graph-batch";

@Injectable({
  providedIn: 'root'
})
export class GraphBatchServiceService {
  private batchUrl = "/api/v1/batch"

  private httpOptions = {
    header: new HttpHeaders({"Content-Type": "application/json"})
  }

  constructor(
    private httpClient: HttpClient
  ) {
  }

  getBatches(): Observable<GraphBatchList> {
    return this.httpClient.get<GraphBatchList>(this.batchUrl)
  }

  getBatch(id: Number): Observable<GraphBatch> {
    const url = `${this.batchUrl}/${id}`
    return this.httpClient.get<GraphBatch>(url)
  }

  createBatch(bat: GraphBatch): Observable<GraphBatch> {
    return this.httpClient.post<GraphBatch>(this.batchUrl, bat);
  }
}
