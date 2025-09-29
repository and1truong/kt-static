# Tìm Hiểu Thánh Kinh

Website giáo dục và Thánh Kinh bằng tiếng Việt, được xây dựng với Docusaurus.

## Mô tả

Đây là một trang web tĩnh dành cho việc nghiên cứu và tìm hiểu Thánh Kinh bằng tiếng Việt. Trang web sử dụng Docusaurus framework để cung cấp trải nghiệm đọc tài liệu tốt với tính năng:

- Tìm kiếm toàn văn
- Bố cục rõ ràng, dễ điều hướng
- Hỗ trợ đa thiết bị (responsive)
- Hệ thống bình luận Disqus
- Tối ưu SEO

## Công nghệ sử dụng

- **Framework**: [Docusaurus 3.9.1](https://docusaurus.io/) - React-based documentation framework
- **Ngôn ngữ**: [TypeScript](https://www.typescriptlang.org/)
- **Frontend**: [React 19.1.1](https://reactjs.org/)
- **Styling**: CSS với [clsx](https://github.com/lukeed/clsx) utility
- **Code Highlighting**: [Prism](https://prismjs.com/)
- **Comments**: [Disqus React](https://github.com/disqus/disqus-react)
- **Deployment**: Vercel

## Cấu trúc thư mục

```
kt-static/
├── VI1934/                 # Nội dung tài liệu
│   ├── [book-codes]/       # Mã sách Thánh Kinh (VU, GI, SU, etc.)
│   │   ├── 1.md           # Chương 1
│   │   ├── 2.md           # Chương 2
│   │   └── _category_.json # Cấu hình danh mục
├── src/
│   ├── css/
│       └── custom.css     # CSS tùy chỉnh
├── static/                # File tĩnh (images, etc.)
├── docusaurus.config.ts   # Cấu hình Docusaurus
├── sidebars.ts           # Cấu hình sidebar
└── package.json          # Dependencies & scripts
```

## Cài đặt và chạy local

### Yêu cầu

- Node.js 22.x (như chỉ định trong package.json)
- npm hoặc yarn

### Các bước thực hiện

1. **Clone repository**
   ```bash
   git clone <repository-url>
   cd kt-static
   ```

2. **Cài đặt dependencies**
   ```bash
   npm install
   ```

3. **Chạy development server**
   ```bash
   npm start
   ```

   Website sẽ chạy tại http://localhost:3000

4. **Build production**
   ```bash
   npm run build
   ```

   File build sẽ được tạo trong thư mục `build/`

### Các lệnh hữu ích

```bash
# Development server
npm start

# Build production
npm run build

# Serve build locally (để test)
npm run serve

# Clear build cache
npm run clear

# Type checking
npm run typecheck
```

## Deployment lên Vercel

### Cách 1: Kết nối GitHub repository (Khuyến khích)

1. Đăng nhập vào [Vercel Dashboard](https://vercel.com/dashboard)
2. Click "New Project"
3. Chọn repository `kt-static` của bạn
4. Vercel sẽ tự động phát hiện Docusaurus project
5. Cấu hình các biến môi trường nếu cần:
   - Không cần biến môi trường đặc biệt cho project này
6. Click "Deploy"

### Cách 2: CLI Deployment

1. **Cài đặt Vercel CLI**
   ```bash
   npm i -g vercel
   ```

2. **Login vào Vercel**
   ```bash
   vercel login
   ```

3. **Deploy project**
   ```bash
   vercel --prod
   ```

### Cấu hình Vercel

Vercel tự động phát hiện Docusaurus project và sử dụng cấu hình sau:

**Build Command**: `npm run build`
**Output Directory**: `build`
**Install Command**: `npm install`

### Environment Variables (nếu cần)

- Không cần environment variables cho project này

## Thêm nội dung mới

1. Tạo thư mục mới trong `VI1934/` với mã sách Thánh Kinh
2. Thêm file `.md` cho từng chương
3. Cập nhật `_category_.json` trong thư mục
4. Cập nhật `sidebars.ts` để thêm vào navigation

## Tùy chỉnh

- **Theme**: Sửa `src/css/custom.css`
- **Navigation**: Sửa `docusaurus.config.ts`
- **Sidebar**: Sửa `sidebars.ts`
- **Cấu hình chung**: Sửa `docusaurus.config.ts`

## License

Copyright © 2024 Tìm Hiểu Thánh Kinh
