import { Link } from "react-router-dom";

import VelocityLogo from "../assets/velocity-logo.svg";

const Velocity = () => (
  <main className="application">
    <img src={VelocityLogo} className="application-logo" alt="Logo of Velocity" />

    <p>
      Edit <code>src/components/Velocity.tsx</code> and save to reload.
    </p>

    <div className="application-links">
      <Link className="application-link" to="/react">
        Go to React page
      </Link>
      <Link className="application-link" to={{ pathname: "https://khulnasoft.io/" }} target="_blank">
        Learn Velocity, a FastHTTP-based Go framework
      </Link>
    </div>
  </main>
);

export default Velocity;
