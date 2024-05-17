import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import Layout from "@theme/Layout";
import HomepageFeatures from "../components/HomepageFeatures";
import { HomepageHeader } from "../components/HomepageHeader";

export default function Home(): JSX.Element {
    const {siteConfig} = useDocusaurusContext();
    
    return (
        <Layout
            title={`${siteConfig.title}`}
            description="Lời Chúa là ngọn đèn cho chân tôi, Ánh sáng cho đường-lối tôi.">
            <HomepageHeader />
            <main>
                <HomepageFeatures />
            </main>
        </Layout>
    );
}
