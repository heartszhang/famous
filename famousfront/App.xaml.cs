using System;
using System.Diagnostics;
using System.Linq;
using System.Threading.Tasks;
using System.Windows;
using famousfront.utils;
using Elysium;
using System.Windows.Media;
namespace famousfront
{
    public partial class App : Application
    {
        protected override void OnStartup(StartupEventArgs e)
        {
            base.OnStartup(e);
            WireUnhandledExceptionHandlers();
            ServiceLocator.Startup();
            this.Apply(Elysium.Theme.Light);
            AlterDefaultBrushes();
        }
        void AlterDefaultBrushes()
        {
          var eg = new Uri("/Elysium;component/Themes/Generic.xaml", UriKind.Relative);
          var lb = new Uri("/Elysium;component/Themes/LightBrushes.xaml", UriKind.Relative);
          var egds = this.Resources.MergedDictionaries.Where(d => d.Source == eg).ToList().SingleOrDefault();
          Debug.Assert(egds != null);
          var brushes = egds.MergedDictionaries.Where(d => d.Source == lb).ToList().SingleOrDefault();
          Debug.Assert(brushes != null);
          brushes["ForegroundBrush"] = ((SolidColorBrush)(new BrushConverter().ConvertFrom("#FF232323"))).GetAsFrozen();
//          brushes["BackgroundBrush"] = ((SolidColorBrush)(new BrushConverter().ConvertFrom("#FFFEFEFE"))).GetAsFrozen();
        }
        static App()
        {
            GalaSoft.MvvmLight.Threading.DispatcherHelper.Initialize();
        }
        protected override void OnExit(ExitEventArgs e)
        {
            ServiceLocator.Shutdown();
            base.OnExit(e);
        }
        [Conditional("RELEASE")]
        private void WireUnhandledExceptionHandlers()
        {
            var h = new UnhandledExceptionEventHandler(CurrentDomain_UnhandledException);
            AppDomain.CurrentDomain.UnhandledException += h.MakeWeakSpecial(x => AppDomain.CurrentDomain.UnhandledException -= x);

            var h2 = new  System.Windows.Threading.DispatcherUnhandledExceptionEventHandler(Dispatcher_UnhandledException);
            Dispatcher.UnhandledException += h2.MakeWeakSpecial(x => Dispatcher.UnhandledException -= x);

            var h3 = new EventHandler<UnobservedTaskExceptionEventArgs>(TaskScheduler_UnobservedTaskException);
            TaskScheduler.UnobservedTaskException += h3.MakeWeak(x => TaskScheduler.UnobservedTaskException -= x);
        }
        static void TaskScheduler_UnobservedTaskException(object sender, UnobservedTaskExceptionEventArgs e)
        {
//          ((App)App.Current).Log.Fatal("Unhandled TaskScheduler Exception", e.Exception);
        }

        static void Dispatcher_UnhandledException(object sender, System.Windows.Threading.DispatcherUnhandledExceptionEventArgs e)
        {
//          ((App)App.Current).Log.Fatal("Unhandled Dispatcher Exception", e.Exception);
        }

        static void CurrentDomain_UnhandledException(object sender, UnhandledExceptionEventArgs e)
        {
          //((App)App.Current).Log.Fatal("Unhandled AppDomain Exception", e.ExceptionObject as Exception);
        }
    }
}
