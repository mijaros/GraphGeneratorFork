import {Component, OnDestroy} from "@angular/core";
import {Subject} from "rxjs";

@Component({
  selector: 'base-component',
  template: ''
})
export class BaseComponent implements OnDestroy {

  protected terminate = new Subject<void>()

  ngOnDestroy(): void {
    this.terminate.next()
    this.terminate.complete()
  }


}
