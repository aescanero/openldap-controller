import * as React from 'react';
import { Admin, Resource, ListGuesser } from 'react-admin';

const Dashboard = () => (
  <Admin>
    <Resource name="users" list={ListGuesser} />
  </Admin>
);

export default Dashboard;
