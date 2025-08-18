import React from 'react';
import { useNavigate } from 'react-router-dom';
import Layout from '../components/Layout';
import LinkButton from '../components/LinkButton';

function HomePage() {
  const navigate = useNavigate();

  return (
    <Layout>
      <div style={{
        display: 'flex',
        position: 'absolute',
        left: 0,
        right: 0,
        top: 0,
        bottom: 0
      }}>
        <div style={{ margin: 'auto' }}>
          <LinkButton onClick={() => navigate('/add-game')}>
            Add game
          </LinkButton>
          
          <LinkButton disabled>
            Player stats
          </LinkButton>
          
          <LinkButton disabled>
            Played games
          </LinkButton>
          
          <LinkButton disabled>
            Game timer
          </LinkButton>
        </div>
      </div>
    </Layout>
  );
}

export default HomePage;
