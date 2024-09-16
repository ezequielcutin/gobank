import React from 'react';
import { Line } from 'react-chartjs-2';
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend } from 'chart.js';
import zoomPlugin from 'chartjs-plugin-zoom';

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, zoomPlugin);

const BalanceChart = ({ balanceHistory }) => {
  const data = {
    labels: balanceHistory.map(item => item.date.toLocaleDateString()),
    datasets: [
      {
        label: 'Balance',
        data: balanceHistory.map(item => item.balance),
        fill: true,
        backgroundColor: 'rgba(99, 102, 241, 0.2)',
        borderColor: 'rgba(99, 102, 241, 1)',
        tension: 0.4,
        pointBackgroundColor: 'rgba(124, 58, 237, 1)',
        pointBorderColor: '#fff',
        pointHoverBackgroundColor: '#fff',
        pointHoverBorderColor: 'rgba(124, 58, 237, 1)'
      }
    ]
  };

  const options = {
    responsive: true,
    plugins: {
      legend: {
        position: 'top',
        labels: {
          font: {
            family: "'Inter', sans-serif",
            size: 12
          },
          color: '#4B5563'
        }
      },
      title: {
        display: true,
        text: 'Balance Over Time',
        font: {
          family: "'Inter', sans-serif",
          size: 16,
          weight: 'bold'
        },
        color: '#1F2937'
      },
      tooltip: {
        backgroundColor: 'rgba(255, 255, 255, 0.8)',
        titleColor: '#1F2937',
        bodyColor: '#4B5563',
        borderColor: 'rgba(99, 102, 241, 0.2)',
        borderWidth: 1,
        bodyFont: {
          family: "'Inter', sans-serif"
        },
        titleFont: {
          family: "'Inter', sans-serif",
          weight: 'bold'
        }
      }
    },
    scales: {
      y: {
        beginAtZero: false,
        grid: {
          color: 'rgba(156, 163, 175, 0.1)'
        },
        ticks: {
          font: {
            family: "'Inter', sans-serif",
            size: 12
          },
          color: '#4B5563'
        }
      },
      x: {
        grid: {
          color: 'rgba(156, 163, 175, 0.1)'
        },
        ticks: {
          font: {
            family: "'Inter', sans-serif",
            size: 12
          },
          color: '#4B5563'
        }
      }
    }
  };

  return (
    <div className="bg-white rounded-xl shadow-lg p-6">
      <Line data={data} options={options} />
    </div>
  );
};

export default BalanceChart;