import NavbarLink from "./NavbarLinks";

export default function Navbar() {
  return (
    <nav className="bg-gray-800 text-white px-6 py-3 shadow-md">
      <div className="flex items-center justify-between max-w-7xl mx-auto">
        <div className="text-xl font-bold">Telemetry Dashboard</div>

        <div className="flex space-x-6">
          <NavbarLink href="/" text="Home"></NavbarLink>
          <NavbarLink href="/trends" text="Trends"></NavbarLink>
          <NavbarLink href="/distributions" text="Distributions"></NavbarLink>
          <NavbarLink href="/kpis" text="KPIs"></NavbarLink>
        </div>
      </div>
    </nav>
  );
}
