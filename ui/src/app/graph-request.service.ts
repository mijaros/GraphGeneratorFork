import {Injectable} from '@angular/core';
import {HttpClient, HttpErrorResponse, HttpHeaders} from '@angular/common/http';
import {Observable, retry, throwError, timer} from "rxjs";
import {GraphRequest} from "./graph-request";
import {GraphRequestList} from "./graph-request-list";

@Injectable({
  providedIn: 'root'
})


export class GraphRequestService {

  private graphsUrl = "/api/v1/graph"
  private httpOptions = {
    header: new HttpHeaders({"Content-Type": "application/json"})
  }

  constructor(
    private httpClient: HttpClient
  ) {
  }

  retryError(error: any, count: Number): Observable<any> {
    if (error instanceof HttpErrorResponse) {
      if (error.status == 503) {
        console.log("Retrying")
        return timer(10 * 1000)
      }
    }
    return throwError(error)
  }

  getGraphs(): Observable<GraphRequestList> {
    return this.httpClient.get<GraphRequestList>(this.graphsUrl).pipe(
      retry({delay: this.retryError})
    )
  }

  getGraph(id: Number): Observable<GraphRequest> {
    const url = `${this.graphsUrl}/${id}`
    return this.httpClient.get<GraphRequest>(url).pipe(
      retry({delay: this.retryError})
    )
  }

  createGraph(graph: GraphRequest): Observable<GraphRequest> {
    return this.httpClient.post<GraphRequest>(this.graphsUrl, graph).pipe(
      retry({delay: this.retryError}))
  }
}
