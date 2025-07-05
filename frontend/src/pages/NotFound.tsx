import { useNavigate } from 'react-router';

const NotFound = () => {
  const navigate = useNavigate();

  return (
    <div className="not-found-page">
      <h1>404 - 页面未找到</h1>
      <p>抱歉，您访问的页面不存在。</p>

      <button 
        onClick={() => navigate('/')}
        className="back-home-button"
      >
        返回首页
      </button>
    </div>
  );
};

export default NotFound;
