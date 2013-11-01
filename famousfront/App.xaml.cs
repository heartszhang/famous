using MetroLog;
using System;
using System.Collections.Generic;
using System.Configuration;
using System.Data;
using System.Diagnostics;
using System.Linq;
using System.Reflection;
using System.Threading.Tasks;
using System.Windows;
using famousfront.core;
namespace famousfront
{
    /// <summary>
    /// Interaction logic for App.xaml
    /// </summary>
    public partial class App : Application
    {
        private static readonly ILogger Log;
        protected override void OnStartup(StartupEventArgs e)
        {
            //Elysium.Manager.Apply(this, Elysium.Theme.Dark, Elysium.AccentBrushes.Blue, System.Windows.Media.Brushes.White);
            WireUnhandledExceptionHandlers();
            base.OnStartup(e);
            var fvi = FileVersionInfo.GetVersionInfo(Assembly.GetExecutingAssembly().Location);
            Log.Info("MishraReader {0} Startup", fvi.ProductVersion);
        }

        static App()
        {
            var fst = new StreamingFileTarget { PathUnderAppData = "famous" };

            LogManagerFactory.DefaultConfiguration.AddTarget(LogLevel.Trace, LogLevel.Fatal, fst);
            Log = LogManagerFactory.DefaultLogManager.GetLogger<App>();
        }
        protected override void OnExit(ExitEventArgs e)
        {
            ServiceLocator.ClearMain();
            ServiceLocator.ClearSettings();
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
            Log.Fatal("Unhandled TaskScheduler Exception", e.Exception);
        }

        static void Dispatcher_UnhandledException(object sender, System.Windows.Threading.DispatcherUnhandledExceptionEventArgs e)
        {
            Log.Fatal("Unhandled Dispatcher Exception", e.Exception);
        }

        static void CurrentDomain_UnhandledException(object sender, UnhandledExceptionEventArgs e)
        {
            Log.Fatal("Unhandled AppDomain Exception", e.ExceptionObject as Exception);
        }
    }
}
