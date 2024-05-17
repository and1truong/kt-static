import clsx from "clsx";
import styles from "../pages/index.module.css";
import Heading from "@theme/Heading";
import useDocusaurusContext from "@docusaurus/core/lib/client/exports/useDocusaurusContext";

export function HomepageHeader() {
    const {siteConfig} = useDocusaurusContext();
    return (
        <header className={clsx('hero hero--primary', styles.heroBanner)}>
            <div className="container">
                <Heading as="h1" className="hero__title">
                    {siteConfig.title}
                </Heading>
                <p className="hero__subtitle">{siteConfig.tagline}</p>
            </div>
        </header>
    );
}
