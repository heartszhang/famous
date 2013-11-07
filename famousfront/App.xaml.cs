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
using Elysium;
namespace famousfront
{
    /// <summary>
    /// Interaction logic for App.xaml
    /// </summary>
    public partial class App : Application
    {

        protected override void OnStartup(StartupEventArgs e)
        {
            this.Apply(Elysium.Theme.Light);
            WireUnhandledExceptionHandlers();
            ServiceLocator.Startup();
            base.OnStartup(e);
            var fvi = FileVersionInfo.GetVersionInfo(Assembly.GetExecutingAssembly().Location);
            ServiceLocator.Log.Info("famousfront {0} Startup", fvi.ProductVersion);
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
            ServiceLocator.Log.Fatal("Unhandled TaskScheduler Exception", e.Exception);
        }

        static void Dispatcher_UnhandledException(object sender, System.Windows.Threading.DispatcherUnhandledExceptionEventArgs e)
        {
            ServiceLocator.Log.Fatal("Unhandled Dispatcher Exception", e.Exception);
        }

        static void CurrentDomain_UnhandledException(object sender, UnhandledExceptionEventArgs e)
        {
            ServiceLocator.Log.Fatal("Unhandled AppDomain Exception", e.ExceptionObject as Exception);
        }
    }
}
