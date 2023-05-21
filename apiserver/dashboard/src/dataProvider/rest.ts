import simpleRestProvider from 'ra-data-simple-rest';
import { fetchUtils } from 'react-admin';

const httpClient = (url: string, options: fetchUtils.Options = {}) => {
    if (!options.headers) {
        options.headers = new Headers({ Accept: 'application/json' });
    }
    const token = localStorage.getItem('token');
    console.log("token: " + token);
    //const authData = token ? JSON.parse(token) : null;
    //const authToken = authData ? authData.token : null;
      if (token) {
          const headers = new Headers(options.headers);
          headers.set('Authorization', `Bearer ${token}`);
          options.headers = headers;
      }
    
    options.method = 'GET';

    return fetchUtils.fetchJson(url, options);
  };

const restProvider = simpleRestProvider('http://localhost:9090', httpClient);

const delayedDataProvider = new Proxy(restProvider, {
    get: (target, name, self) =>
        name === 'then' // as we await for the dataProvider, JS calls then on it. We must trap that call or else the dataProvider will be called with the then method
            ? self
            : (resource: string, params: any) =>
                  new Promise(resolve =>
                      setTimeout(
                          () =>
                              resolve(
                                  restProvider[name as string](resource, params)
                              ),
                          500
                      )
                  ),
});

export default delayedDataProvider;
