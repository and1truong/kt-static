import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  Svg: React.ComponentType<React.ComponentProps<'svg'>>;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Sự mặc khải từ thiên thượng',
    Svg: require('@site/static/img/undraw_docusaurus_mountain.svg').default,
    description: (
      <>
        Kinh Thánh là ngọn đèn cho chân chúng ta và là ánh sáng cho đường đi của chúng ta.
        Lời của Thiên Chúa dẫn dắt chúng ta trong suốt cuộc đời -- 
        cung cấp cho chúng ta các nguyên tắc để tuân theo để giúp chúng ta trong việc ra quyết định.
      </>
    ),
  },
  {
    title: 'Nguyên liệu dưỡng linh quan trọng nhất',
    Svg: require('@site/static/img/undraw_docusaurus_tree.svg').default,
    description: (
      <>
        Kinh Thánh là nguồn kiến thức và sự khôn ngoan tối thượng.
        Học hỏi kinh văn mỗi ngày góp phần xây dựng đức tin của mình nơi Thiên Chúa.
      </>
    ),
  },
  {
    title: 'Nguồn ơn phước vô tận',
    Svg: require('@site/static/img/undraw_docusaurus_react.svg').default,
    description: (
      <>
        Đọc Kinh Thánh là khoản đầu tư tốt nhất trong cuộc sống.
        Nếu chúng ta giữ mọi điều đã viết trong luật pháp, chúng ta sẽ thành công và thịnh vượng.
      </>
    ),
  },
];

function Feature({title, Svg, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} role="img" />
      </div>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): JSX.Element {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
