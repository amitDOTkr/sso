import { Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor
} from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable()
export class ApiInterceptorInterceptor implements HttpInterceptor {

  constructor() {}

  intercept(request: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
    const baseUrl = document.getElementsByTagName('base')[0].href;
    // console.log(baseUrl)
    if (!request.url.startsWith("http://") && !request.url.startsWith("https://")) {
      const apiReq = request.clone({ url: `${baseUrl}${request.url}` });
      // console.log(apiReq)
      return next.handle(apiReq);
    }
    return next.handle(request)
  }
}
