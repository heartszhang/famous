using famousfront.viewmodels;
using System;
using System.Windows.Media;
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
      // expression windowsinstance to window error if IsMainWindow setted in xaml
      SetIsMainWindow(this, true);
      Elysium.Parameters.General.SetShadowBrush(this, Brushes.Black);
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
