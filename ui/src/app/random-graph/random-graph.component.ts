import {Component, OnInit} from '@angular/core';
import {GraphRequest} from "../graph-request";
import {GraphRequestService} from "../graph-request.service";
import {explanation} from "../graph-type";
import {MessageService} from "../message.service";
import {Severity} from "../message";
import {faCircleInfo, faDownload, faPlus} from "@fortawesome/free-solid-svg-icons";
import {HttpErrorResponse} from "@angular/common/http";
import {KeyValue} from "@angular/common";
import {explain, getClass, RequestStatus} from "../request-status";
import {takeUntil} from "rxjs";
import {BaseComponent} from "../base.component";
import {DateTime} from "luxon";

@Component({
  selector: 'app-random-graph',
  templateUrl: './random-graph.component.html',
  styleUrls: ['./random-graph.component.css']
})
export class RandomGraphComponent extends BaseComponent implements OnInit {
  readonly pageSize = 6
  graphs: Map<number, GraphRequest> = new Map<number, GraphRequest>()
  graphOrder: Number[] = []
  showForm: boolean = false
  expandedDetails: number | null = null
  page: number = 1
  protected readonly explanation = explanation;
  protected readonly faDownload = faDownload;
  protected readonly explain = explain;
  protected readonly getClass = getClass;
  protected readonly RequestStatus = RequestStatus;
  protected readonly faPlus = faPlus;
  protected readonly faCircleInfo = faCircleInfo;
  protected readonly Intl = Intl;
  private interval: any
  private allRefresh: any

  constructor(private graphRequestService: GraphRequestService,
              private messageService: MessageService) {
    super()
  }

  get requests(): GraphRequest[] {
    return this.graphOrder.map(value => this.graphs.get(value.valueOf())!)
  }

  getGraphs(): void {
    this.graphRequestService.getGraphs().pipe(takeUntil(this.terminate)).subscribe(graphList => {
      graphList.graphs.forEach(i => {
        this.getGraph(i)
      })
      let eliminated: number[] = []
      for (let k of this.graphs.keys()) {
        if (!graphList.graphs.indexOf(k)) {
          eliminated.push(k)
        }
      }
      eliminated.forEach(value => {
        this.graphs.delete(value)
      })
    })
  }

  submitForm(graph: GraphRequest | null): void {
    this.showForm = false
    if (graph != null) {
      this.addGraph(graph!)
    }

  }

  graphInDeletedOrder = (a: KeyValue<number, GraphRequest>, b: KeyValue<number, GraphRequest>): number => {
    return b.value.deleted.valueOf() > a.value.deleted.valueOf() ? -1 : 1
  }

  addGraph(graph: GraphRequest): void {
    this.messageService.addMessage(`Graph was requested: id=${graph.id}`, Severity.SUCCESS)
    this.graphs.set(graph.id.valueOf(), graph)
    this.graphOrder.push(graph.id)
  }

  getGraph(id: Number): void {
    if (this.graphs.has(id.valueOf())) {
      let graph = this.graphs.get(id.valueOf())!
      if (graph.status == RequestStatus.DELETED ||
        (graph.status == RequestStatus.FINISHED && DateTime.fromJSDate(graph.deleted).diff(DateTime.now(), 'seconds').get("seconds") > 60)) {
        return
      }
    }
    this.graphRequestService.getGraph(id).pipe(takeUntil(this.terminate)).subscribe({
      next: (gr) => {
        if (!this.graphs.has(id.valueOf())) {
          gr.deleted = new Date(gr.deleted)
          this.graphOrder.push(id)
        }
        this.graphs.set(gr.id.valueOf(), gr)
      }, error: (err) => {
        if (err instanceof HttpErrorResponse) {
          if (err.status == 404) {
            let graph = this.graphs.get(id.valueOf())!
            graph.status = RequestStatus.DELETED
            this.graphs.set(id.valueOf(), graph)
          }
        }
      }
    })
  }

  updateGraphs(): void {
    this.graphs.forEach((_, id) => this.getGraph(id))
  }

  //
  ngOnInit(): void {
    this.getGraphs()
    this.interval = setInterval(() => {
      this.updateGraphs()
    }, 30 * 1000)
    this.allRefresh = setInterval(() => {
      this.getGraphs()
    }, 3 * 60 * 1000)
  }

  override ngOnDestroy() {
    clearInterval(this.interval)
    clearInterval(this.allRefresh)
    super.ngOnDestroy();
  }

  toggleForm(): void {
    this.showForm = !this.showForm
  }

  detailToggle(id: number) {
    if (id == this.expandedDetails) {
      this.expandedDetails = null
    } else {
      this.expandedDetails = id
    }
  }
}
