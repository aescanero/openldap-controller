import React from 'react';
import { Admin, Resource } from 'react-admin';
import { Graph } from './graph';
import { Login, Layout } from './layout';
import delayedDataProvider from './dataProvider';
import authProvider from './authProvider';



const App = () => (
  <Admin
        title=""
        dataProvider={delayedDataProvider}
        authProvider={authProvider}
        dashboard={Graph}
        loginPage={Login}
        layout={Layout}
        //i18nProvider={i18nProvider}
        disableTelemetry
        //theme={lightTheme}
        requireAuth
    >
    <Resource name="graphs" {...Graph} />
  </Admin>
);  

export default App;
