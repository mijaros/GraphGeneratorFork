import {Component, OnInit} from '@angular/core';
import {GraphBatchServiceService} from "../graph-batch-service.service";
import {GraphBatch} from "../graph-batch";
import {explanation} from "../graph-type";
import {MessageService} from "../message.service";
import {Severity} from "../message";
import {faCircleInfo, faDownload} from "@fortawesome/free-solid-svg-icons";
import {explain, getClass, RequestStatus} from "../request-status";
import {HttpErrorResponse} from "@angular/common/http";
import {Subject, takeUntil} from "rxjs";
import {BaseComponent} from "../base.component";
import {KeyValue} from "@angular/common";
import {DateTime} from "luxon"

@Component({
  selector: 'app-random-batch',
  templateUrl: './random-batch.component.html',
  styleUrls: ['./random-batch.component.css']
})

export class RandomBatchComponent extends BaseComponent implements OnInit {
  batches: Map<Number, GraphBatch> = new Map<Number, GraphBatch>()
  readonly pageSize = 6
  page = 1
  showFrom: boolean = false
  detailsToggle: number | null = null
  subs = new Subject<void>()
  protected readonly explanation = explanation;
  protected readonly faDownload = faDownload;
  protected readonly explain = explain;
  protected readonly getClass = getClass;
  protected readonly RequestStatus = RequestStatus;
  protected readonly faCircleInfo = faCircleInfo;
  protected readonly Intl = Intl;
  private interval: any
  private individual: any

  constructor(private batchService: GraphBatchServiceService, private messageService: MessageService) {
    super()
  }

  batchInDeletedOrder = (a: KeyValue<Number, GraphBatch>, b: KeyValue<Number, GraphBatch>): number => {
    return b.value.deleted.valueOf() > a.value.deleted.valueOf() ? -1 : 1
  }

  formToggle() {
    this.showFrom = !this.showFrom
  }

  getBatches(): void {
    this.batchService.getBatches().pipe(takeUntil(this.terminate)).subscribe(bat => {
      let rem: number[] = []
      bat.batches.forEach(bt => {
        if (!this.batches.has(bt) || this.batches.get(bt)!.status === RequestStatus.IN_PROGRESS) {
          this.getBatch(bt)
        }
      })
      for (let id of this.batches.keys()) {
        if (bat.batches.indexOf(id)) {
          continue
        }
        rem.push(id.valueOf())
      }
      rem.forEach(value => this.batches.delete(value))
    })
  }

  handleForm(bat: GraphBatch | null) {
    this.showFrom = false
    if (bat != null) {
      this.addBatch(bat!)
    }
  }

  getBatch(id: Number): void {
    if (this.batches.has(id)) {
      let batch = this.batches.get(id)!
      if ((batch.status == RequestStatus.FINISHED && DateTime.fromJSDate(batch.deleted).diff(DateTime.now(), 'seconds').get('seconds') > 60)
        || batch.status == RequestStatus.DELETED
      ) {
        return
      }
    }
    this.batchService.getBatch(id).pipe(takeUntil(this.terminate)).subscribe({
      next: value => {
        value.deleted = new Date(value.deleted)
        this.batches.set(id, value)
      },
      error: err => {
        if (err instanceof HttpErrorResponse && err.status == 404) {
          console.log(`request ${id} was deleted`)
          let batch = this.batches.get(id)!
          batch.status = RequestStatus.DELETED
          this.batches.set(id, batch)
        } else if (err instanceof HttpErrorResponse && err.status == 503) {
          this.messageService
        }
      }
    })
  }

  addBatch(bat: GraphBatch): void {
    this.messageService.addMessage("Batch successfully created.", Severity.SUCCESS)
    this.batches.set(bat!.id, bat!)
  }

  ngOnInit(): void {
    this.getBatches()
    this.interval = setInterval(() => {
      this.getBatches()
    }, 5 * 60 * 1000)
    this.individual = setInterval(() => {
      this.batches.forEach((value, key) => {
        this.getBatch(key)
      })
    }, 10 * 1000)
  }

  override ngOnDestroy() {
    if (this.individual) {
      clearInterval(this.individual)
    }
    if (this.interval) {
      clearInterval(this.interval)
    }
    super.ngOnDestroy()
  }

  toggleDetails(n: number) {
    if (this.detailsToggle == n) {
      this.detailsToggle = null
    } else {
      this.detailsToggle = n
    }
  }

  graphIdList(ids: Number[]): string {
    return ids.join(", ")
  }
}
