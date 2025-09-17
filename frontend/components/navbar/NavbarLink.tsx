import Link from "next/link";

type NavbarLinkProps = {
  href: string;
  text: string;
};

export default function NavbarLink(props: NavbarLinkProps) {
  return (
    <Link href={props.href} className="hover:text-blue-400 transition">
      {props.text}
    </Link>
  );
}
