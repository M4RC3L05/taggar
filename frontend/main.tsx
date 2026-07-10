// @ts-ignore
import "./main.css"
import { render } from "preact";
import { Router, Route, Switch } from "wouter-preact";
import { useHashLocation } from "wouter-preact/use-hash-location";
import Home from "./src/pages/home";
import ManageFile from "./src/pages/manage-file";

const App = () => {
  return <>
    <Router hook={useHashLocation}>
      <Switch>
        <Route path="/" component={Home} />
        <Route path="/manage-file/:file" component={ManageFile} />
      </Switch>
    </Router>
  </>
}

render(<App />, document.getElementById("app")!)
