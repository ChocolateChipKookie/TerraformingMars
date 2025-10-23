import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Layout from '../components/Layout';
import MetalContainer from '../components/MetalContainer';

function PlayerStatsPage() {
  const navigate = useNavigate();
  const [playerStats, setPlayerStats] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch('/data/player_stats.json')
      .then(response => response.json())
      .then(data => {
        setPlayerStats(data);
        setLoading(false);
      })
      .catch(error => {
        console.error('Error loading player stats:', error);
        setLoading(false);
      });
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <div>Player stats</div>

      <div>
        <div>Statistics</div>
        <div>
          <div>
            <table>
              <thead>
                <tr>
                  <th>Player</th>
                  <th>Games</th>
                  <th>Wins</th>
                  <th>Win %</th>
                  <th>Avg Points</th>
                  <th>ELO</th>
                </tr>
              </thead>
              <tbody>
                {playerStats && Object.entries(playerStats).map(([player, stats]) => (
                  <tr key={player}>
                    <td>{player}</td>
                    <td>{stats.games}</td>
                    <td>{stats.wins}</td>
                    <td>{((stats.wins / stats.games) * 100).toFixed(1)}%</td>
                    <td>{stats.avgPoints?.toFixed(1) || '-'}</td>
                    <td>{stats.elo || '-'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <div>
        <button onClick={() => navigate('/')}>Main page</button>
      </div>
    </div>
  );
}

export default PlayerStatsPage;