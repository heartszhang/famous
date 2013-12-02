using famousfront.viewmodels;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Navigation;
using System.Windows.Shapes;
//using Elysium.Parameters;
namespace famousfront
{
    /// <summary>
    /// Interaction logic for MainWindow.xaml
    /// </summary>
    public partial class MainWindow : Elysium.Controls.Window
    {
        MainViewModel _main_viewmodel = null;
        public MainWindow()
        {          
            InitializeComponent();            
            SetIsMainWindow(this, true);
            Elysium.Parameters.General.SetShadowBrush(this, Brushes.Black);
            // expression windowsinstance to window error if IsMainWindow setted in xaml
        }
        protected override void OnInitialized(EventArgs e)
        {
            base.OnInitialized(e);
            var locator = DataContext as ServiceLocator;
            if (locator == null)
              return;
            _main_viewmodel = locator.MainViewModel;
            _content.Content = _main_viewmodel;
        }
    }
}
