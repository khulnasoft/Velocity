import { BrowserRouter, Route, Switch } from "react-router-dom";

import Velocity from "./components/Velocity";
import NotFound from "./components/NotFound";
import React from "./components/React";

const App = () => (
  // Add basename to the <BrowserRouter basename="/web"> if you serve Single Page Application on "/web"
  <BrowserRouter>
    <Switch>
      <Route path="/" component={Velocity} exact />
      <Route path="/react" component={React} exact />
      <Route path="*" component={NotFound} />
    </Switch>
  </BrowserRouter>
);

export default App;
