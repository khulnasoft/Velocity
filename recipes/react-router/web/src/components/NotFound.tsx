import { Link } from "react-router-dom";

import VelocityLogo from "../assets/velocity-logo.svg";
import ReactLogo from "../assets/react-logo.svg";

const NotFound = () => (
  <main className="application">
    <img src={VelocityLogo} className="application-logo" alt="Logo of Velocity" />
    <img src={ReactLogo} className="application-logo" alt="Logo of React" />

    <p>Page not found! Let's go back home!</p>

    <div className="application-links">
      <Link className="application-link" to="/">
        Back home
      </Link>
    </div>
  </main>
);

export default NotFound;
